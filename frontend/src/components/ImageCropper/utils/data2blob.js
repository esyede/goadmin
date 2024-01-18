/**
 * Convert base64 file format to binary
 *
 * @param  {[String]} data data:image/png;base64,[this part]xxxxxx[/this part]
 * @param  {[String]} mime mime-type
 * @return {[blob]}
 */
export default function (data, mime) {
  data = data.split(',')[1]
  data = window.atob(data)
  let ia = new Uint8Array(data.length)
  for (let i = 0; i < data.length; i++) {
    ia[i] = data.charCodeAt(i)
  }
  return new Blob([ia], { type: mime })
}
