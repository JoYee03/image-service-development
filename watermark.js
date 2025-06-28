const sharp = require('sharp');

// Input files
const IMAGE_PATH = 'temp-image.jpg';
const WATERMARK_PATH = 'temp-watermark.jpg';
const OUTPUT_PATH = 'output.jpg';

(async () => {
  try {
    const image = sharp(IMAGE_PATH);
    const watermark = await sharp(WATERMARK_PATH)
    .resize({ width: 100, height: 100 })
    .toBuffer();

    const { width: imgWidth, height: imgHeight } = await image.metadata();
    const { width: wmWidth, height: wmHeight } = await sharp(watermark).metadata();

    if (!imgWidth || !imgHeight || !wmWidth || !wmHeight) {
      throw new Error('Error reading image or watermark size.');
    }

    // Repeat watermark across image
    const compositeArray = [];
    for (let y = 0; y < imgHeight; y += wmHeight) {
      for (let x = 0; x < imgWidth; x += wmWidth) {
        compositeArray.push({ input: watermark, left: x, top: y });
      }
    }

    await image
      .composite(compositeArray)
      .toFile(OUTPUT_PATH);

    console.log('Watermarking complete.');
  } catch (error) {
    console.error('Error in watermark.js:', error.message);
    process.exit(1);
  }
})();
