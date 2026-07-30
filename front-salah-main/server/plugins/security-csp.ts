import { randomBytes } from 'node:crypto';
import { setResponseHeader } from 'h3';

const scriptWithoutNonce = /<script(?![^>]*\bnonce=)/gi;
const htmlSections = ['head', 'bodyPrepend', 'body', 'bodyAppend'] as const;

export default defineNitroPlugin((nitroApp) => {
  nitroApp.hooks.hook('render:html', (html, { event }) => {
    const nonce = randomBytes(18).toString('base64');

    for (const section of htmlSections) {
      html[section] = html[section].map((chunk) =>
        chunk.replace(scriptWithoutNonce, `<script nonce="${nonce}"`),
      );
    }

    setResponseHeader(event, 'Content-Security-Policy', [
      "default-src 'self'",
      "img-src 'self' data: blob:",
      "style-src 'self' 'unsafe-inline'",
      `script-src 'self' 'nonce-${nonce}'`,
      "script-src-attr 'none'",
      "connect-src 'self'",
      "object-src 'none'",
      "base-uri 'self'",
      "frame-ancestors 'none'",
      "form-action 'self'",
    ].join('; '));
  });
});
