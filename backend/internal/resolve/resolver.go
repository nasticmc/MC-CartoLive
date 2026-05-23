package resolve

import (
	"context"
	"fmt"

	"meshcore-australia-live-map/backend/internal/meshcore"
)

type CandidateProvider interface {
	CandidatesByPrefix(ctx context.Context, iata string, hashSize int, prefix string) ([]Candidate, error)
}

type Resolver struct {
	provider       CandidateProvider
	forwarderRoles map[string]bool
}

func New(provider CandidateProvider, forwarderRoles []string) *Resolver {
	roles := map[string]bool{}
	for _, role := range forwarderRoles {
		roles[role] = true
	}
	if len(roles) == 0 {
		roles["repeater"] = true
		roles["room_server"] = true
	}
	return &Resolver{provider: provider, forwarderRoles: roles}
}

func (r *Resolver) Resolve(ctx context.Context, iata string, parsed meshcore.ParsedPacket) (Result, error) {
	if parsed.InvalidForMap {
		return Result{Status: StatusInvalidForMap, Reason: parsed.InvalidReason}, nil
	}
	if parsed.HopCount == 0 {
		return Result{Status: StatusNoPath, Reason: "zero_hop_packet"}, nil
	}

	seen := map[string]bool{}
	for _, prefix := range parsed.PathChunks {
		if seen[prefix] {
			return Result{Status: StatusDuplicatePrefix, Reason: fmt.Sprintf("duplicate prefix %s", prefix)}, nil
		}
		seen[prefix] = true
	}

	result := Result{Status: StatusHigh}
	for _, prefix := range parsed.PathChunks {
		candidates, err := r.provider.CandidatesByPrefix(ctx, iata, parsed.HashSize, prefix)
		if err != nil {
			return Result{}, err
		}
		if len(candidates) == 0 {
			return Result{Status: StatusUnresolved, Reason: fmt.Sprintf("no candidates for %d-byte prefix %s in %s", parsed.HashSize, prefix, iata)}, nil
		}

		forwarders := make([]Candidate, 0, len(candidates))
		for _, candidate := range candidates {
			if r.forwarderRoles[candidate.Role] {
				forwarders = append(forwarders, candidate)
			}
		}
		if len(forwarders) == 0 {
			return Result{Status: StatusRoleInvalid, Reason: fmt.Sprintf("prefix %s has no forwarder-capable candidates", prefix)}, nil
		}
		if len(forwarders) > 1 {
			return Result{Status: StatusAmbiguous, Reason: fmt.Sprintf("prefix %s maps to %d forwarder candidates", prefix, len(forwarders))}, nil
		}

		result.Hops = append(result.Hops, ResolvedHop{
			Prefix:     prefix,
			Confidence: StatusHigh,
			Candidate:  forwarders[0],
		})
	}

	return result, nil
}
