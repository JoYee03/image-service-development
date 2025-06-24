const functions = require('firebase-functions')
const { execFile } = require('child_process')

exports.testImageUpload = functions.https.onRequest((req, res) => {
  execFile(
    './image-service',             // Path to your Go binary
    ['--upload', req.body.content, req.body.type], // or whatever flags you implement
    (error, stdout, stderr) => {
      if (error) {
        console.error(error, stderr)
        return res.status(500).send({ success: false, error: error.message })
      }
      const response = JSON.parse(stdout)
      return res.status(200).send(response)
    }
  )
})

exports.testWatermarkImage = functions.https.onRequest((req, res) => {
  execFile(
    './image-service',
    ['--watermark', req.body.image_path, req.body.watermark_path],
    (error, stdout, stderr) => {
      if (error) {
        console.error(error, stderr)
        return res.status(500).send({ success: false, error: error.message })
      }
      const response = JSON.parse(stdout)
      return res.status(200).send(response)
    }
  )
})
