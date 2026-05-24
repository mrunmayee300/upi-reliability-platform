import { readFileSync, writeFileSync } from "fs";
import { Resvg } from "@resvg/resvg-js";
import { resolve, dirname } from "path";
import { fileURLToPath } from "url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const imagesDir = resolve(__dirname, "../../docs/images");

const files = ["architecture-overview.svg", "transaction-flow.svg"];

for (const svgFile of files) {
  const svgPath = resolve(imagesDir, svgFile);
  const pngPath = svgPath.replace(/\.svg$/i, ".png");
  const svg = readFileSync(svgPath, "utf8");
  const resvg = new Resvg(svg, {
    fitTo: { mode: "width", value: 1600 },
  });
  const png = resvg.render().asPng();
  writeFileSync(pngPath, png);
  console.log(`Wrote ${pngPath} (${png.length} bytes)`);
}
