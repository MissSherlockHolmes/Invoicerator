/**
 * PDF.js loaded as ESM (no global pdfjsLib). Cached after first use.
 * v4 matches worker + main bundle from the same release.
 */
const PDFJS_VERSION = '4.8.69';
const PDF_MJS = `https://cdn.jsdelivr.net/npm/pdfjs-dist@${PDFJS_VERSION}/build/pdf.mjs`;
const PDF_WORKER = `https://cdn.jsdelivr.net/npm/pdfjs-dist@${PDFJS_VERSION}/build/pdf.worker.mjs`;

let pdfApiPromise = null;

export async function loadPdfApi() {
    if (!pdfApiPromise) {
        pdfApiPromise = import(PDF_MJS).then((m) => {
            if (m.GlobalWorkerOptions) {
                m.GlobalWorkerOptions.workerSrc = PDF_WORKER;
            }
            return m;
        });
    }
    return pdfApiPromise;
}

export async function extractTextFromPDF(arrayBuffer) {
    const { getDocument } = await loadPdfApi();
    const data = arrayBuffer instanceof ArrayBuffer ? new Uint8Array(arrayBuffer) : arrayBuffer;
    const pdf = await getDocument({ data }).promise;
    let fullText = '';

    for (let i = 1; i <= pdf.numPages; i++) {
        const page = await pdf.getPage(i);
        const textContent = await page.getTextContent();

        let text = '';
        for (const item of textContent.items) {
            text += item.str;
            if (item.hasEOL) {
                text += '\n';
            } else {
                text += ' ';
            }
        }

        fullText += text + '\n';
    }

    return fullText;
}
