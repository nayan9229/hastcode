// Batch Mermaid renderer, run by internal/mermaid as
//   node --input-type=module --eval <this file>   (cwd = the tools dir)
// stdin: JSON array of Mermaid sources. stdout: JSON array of
// {svg} | {error}, one per source, in order.
import { renderMermaidSVG } from 'beautiful-mermaid';

let input = '';
for await (const chunk of process.stdin) input += chunk;

const out = JSON.parse(input).map((src) => {
  try {
    return {
      svg: renderMermaidSVG(src, {
        // Theme hooks resolved by the site CSS (.diagram in diagram.css).
        // Do not name these --bg/--accent: the renderer emits
        // --accent:var(--accent), a cycle.
        bg: 'var(--diagram-bg)',
        fg: 'currentColor',
        accent: 'var(--diagram-accent)',
        transparent: true,
      }),
    };
  } catch (e) {
    return { error: String(e && e.message ? e.message : e) };
  }
});
process.stdout.write(JSON.stringify(out));
