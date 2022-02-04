import str from '../package.json?raw'

export function name() {
  console.log(import.meta.env.NODE_ENV)
  return JSON.parse(str).name
}
