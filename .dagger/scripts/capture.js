// Captures the Svelte presentation as a sequence of PNG screenshots,
// preceded by a title page that links to the GitHub repo.
//
// Invoked by the ExportPdf Dagger function. Expects:
//   BASE_URL   — where the presentation is served (default http://presentation/)
//   REPO_URL   — repo link baked into the title page
//   OUT_DIR    — directory for PNG output (default /work/out)

const { chromium } = require('playwright');
const fs = require('fs');
const path = require('path');

const REPO_URL = process.env.REPO_URL || 'https://github.com/AsocPro/bsides-k8s-2026';
const BASE_URL = process.env.BASE_URL || 'http://presentation/';
const OUT_DIR  = process.env.OUT_DIR  || '/work/out';

// Slide order mirrors frontend/src/App.svelte. `substeps` is the number of
// ArrowRight presses that reveal progressive content on the slide before the
// engine advances to the next one (driven by registerSubsteps()).
const SLIDES = [
  { id: 'title',        substeps: 0 },
  { id: 'koolaid',      substeps: 0 },
  { id: 'evolution',    substeps: 3 },
  { id: 'k8s-intro',    substeps: 4 },
  { id: 'k8s-security', substeps: 0 },
  { id: 'rbac',         substeps: 0 },
  { id: 'rbac-demo',    substeps: 0 },
  { id: 'policy',       substeps: 0 },
  { id: 'policy-demo',  substeps: 0 },
  { id: 'netpol',       substeps: 0 },
  { id: 'netpol-demo',  substeps: 0 },
  { id: 'future',       substeps: 0 },
  { id: 'end',          substeps: 0 },
];

const VIEWPORT = { width: 1920, height: 1080 };

function titleHTML(repoUrl) {
  return `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  html, body {
    height: 100%;
    background: #0a0a0a;
    color: #e5e5e5;
    font-family: 'Inter', system-ui, -apple-system, sans-serif;
  }
  .page {
    width: 100vw; height: 100vh;
    display: flex; flex-direction: column;
    justify-content: center; align-items: center;
    padding: 4rem; text-align: center; position: relative;
  }
  .eyebrow {
    font-size: 1rem; color: #737373;
    text-transform: uppercase; letter-spacing: 0.25em;
    margin-bottom: 2.5rem;
  }
  .title {
    font-size: 4.25rem; font-weight: 700;
    margin-bottom: 1.25rem;
    background: linear-gradient(135deg, #22d3ee, #a78bfa);
    -webkit-background-clip: text; background-clip: text;
    color: transparent;
    max-width: 85%; line-height: 1.08;
  }
  .byline {
    font-size: 1.5rem; color: #a3a3a3;
    margin-bottom: 4rem;
  }
  .repo-label {
    font-size: 0.875rem; color: #737373;
    text-transform: uppercase; letter-spacing: 0.18em;
    margin-bottom: 0.9rem;
  }
  .repo-link {
    font-family: ui-monospace, 'SF Mono', Menlo, Consolas, monospace;
    font-size: 1.5rem; color: #22d3ee;
    padding: 0.9rem 1.75rem;
    border: 1px solid rgba(34, 211, 238, 0.25);
    border-radius: 0.5rem;
    background: rgba(34, 211, 238, 0.05);
    text-decoration: none;
  }
  .footer {
    position: absolute; bottom: 3rem;
    font-size: 0.95rem; color: #525252;
  }
  .footer code {
    font-family: ui-monospace, 'SF Mono', Menlo, Consolas, monospace;
    background: #1f1f1f; padding: 0.15rem 0.5rem;
    border-radius: 0.3rem; color: #cbd5e1;
  }
</style>
</head>
<body>
<div class="page">
  <p class="eyebrow">BSides &middot; Zach Gibbs</p>
  <h1 class="title">Kubernetes<br>Orchestrating a More Secure Future</h1>
  <p class="byline">Live, interactive demos with real k3s clusters</p>
  <p class="repo-label">Source, setup, and full demo environment</p>
  <a class="repo-link" href="${repoUrl}">${repoUrl}</a>
  <p class="footer">Clone the repo and run <code>dagger call present up --ports 8080:8080</code> to replay everything yourself.</p>
</div>
</body>
</html>`;
}

(async () => {
  fs.mkdirSync(OUT_DIR, { recursive: true });

  // Write the title page next to the script so Playwright can load it via file://
  const titlePath = path.join(OUT_DIR, '_title.html');
  fs.writeFileSync(titlePath, titleHTML(REPO_URL));

  const browser = await chromium.launch();
  const context = await browser.newContext({ viewport: VIEWPORT });
  const page = await context.newPage();

  // 1. Title page — rendered as PDF (not PNG) so the <a> tag becomes a real,
  // clickable PDF link annotation. Page dimensions match the 1920x1080 PNG
  // screenshots at 96 DPI (20in x 11.25in) so pdfunite produces a deck with
  // a consistent page size.
  await page.goto('file://' + titlePath);
  await page.waitForLoadState('networkidle').catch(() => {});
  // Stay in "screen" media so our dark-mode styles apply during print.
  await page.emulateMedia({ media: 'screen' });
  await page.evaluate(() => document.fonts && document.fonts.ready).catch(() => {});
  await page.pdf({
    path: path.join(OUT_DIR, 'title.pdf'),
    width: '20in',
    height: '11.25in',
    printBackground: true,
    margin: { top: 0, right: 0, bottom: 0, left: 0 },
    pageRanges: '1',
  });
  console.log('Captured title page -> title.pdf (link is clickable)');

  // 2. Presentation — step through each slide, triggering any substeps
  await page.goto(BASE_URL, { waitUntil: 'domcontentloaded' });
  await page.waitForSelector('.slide', { timeout: 30000 });
  // Let Svelte/GSAP settle on the first slide
  await page.waitForTimeout(2000);

  for (let i = 0; i < SLIDES.length; i++) {
    const slide = SLIDES[i];

    // Reveal substeps (GSAP animations need time to complete)
    for (let s = 0; s < slide.substeps; s++) {
      await page.keyboard.press('ArrowRight');
      await page.waitForTimeout(900);
    }
    // Final settle before capture
    await page.waitForTimeout(600);

    const pageNum = String(i + 1).padStart(2, '0');
    const outPath = path.join(OUT_DIR, `page-${pageNum}.png`);
    await page.screenshot({ path: outPath, type: 'png' });
    console.log(`Captured slide ${i + 1}/${SLIDES.length}: ${slide.id} -> page-${pageNum}.png`);

    // Advance to the next slide (unless we're on the last one)
    if (i < SLIDES.length - 1) {
      await page.keyboard.press('ArrowRight');
      await page.waitForTimeout(800);
    }
  }

  await browser.close();
})();
