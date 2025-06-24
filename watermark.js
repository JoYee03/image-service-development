const sharp = require('sharp');

(async () => {
  try {
    await sharp('temp-image.png')
      .composite([{
        input: 'temp-watermark.png',
        gravity: 'southeast'
      }])
      .toFile('output.png');

    console.log("Watermark applied successfully");
  } catch (err) {
    console.error("Sharp error:", err);
    process.exit(1);
  }
})();
