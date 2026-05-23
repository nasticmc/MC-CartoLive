import { ExternalLink, Github } from 'lucide-react';
import { appVersion, buildNumber } from '../buildInfo';

const GITHUB_URL = 'https://github.com/n30nex/MC-CartoLive';
const EASTMESH_AU_URL = 'https://eastmesh.au/';

export default function LinkBar() {
  return (
    <nav className="link-bar" aria-label="Project links">
      <a className="link-bar-brand" href={EASTMESH_AU_URL} target="_blank" rel="noreferrer" title="Open EastMesh">
        <img src="/meshcore-australia-favicon.png" alt="" aria-hidden="true" />
        <span>EastMesh</span>
      </a>
      <div className="link-bar-build" aria-label={`MC-CartoLive version ${appVersion}, build ${buildNumber}`}>
        <strong>MC-CartoLive</strong>
        <span>v{appVersion}</span>
        <span>build {buildNumber}</span>
      </div>
      <a className="link-bar-github" href={GITHUB_URL} target="_blank" rel="noreferrer" title="Open MC-CartoLive on GitHub">
        <Github size={15} />
        <span>GitHub</span>
        <ExternalLink size={12} />
      </a>
    </nav>
  );
}
