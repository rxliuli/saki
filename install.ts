import * as os from 'os'
import json from './package.json'
import { downloadRelease } from '@terascope/fetch-github-release'
import * as path from 'path'
import { extract } from 'tar'
import { copy } from 'fs-extra'

async function main() {
  const archMap: Record<string, string> = {
    arm64: 'arm64',
    x64: 'amd64',
  }
  const platformMap: Record<string, string> = {
    win32: 'windows',
    linux: 'linux',
    darwin: 'macos',
  }
  const platform = platformMap[os.platform()]
  const arch = os.arch()
  const assetName = `saki_${json.version}_${platform}_${archMap[arch]}.tar.gz`
  const tempPath = path.resolve(__dirname, '.temp')
  try {
    await downloadRelease(
      'rxliuli',
      'saki',
      tempPath,
      (release) => {
        return release.tag_name === `v${json.version}`
      },
      (asset) => {
        return asset.name === assetName
      },
      false,
      false,
    )
  } catch (e) {
    console.error('下载失败：', e)
  }
  await extract({
    file: path.resolve(tempPath, assetName),
    cwd: tempPath,
  })
  await copy(
    path.resolve(tempPath, 'saki' + (platform === 'windows' ? '.exe' : '')),
    path.resolve(__dirname, 'bin'),
  )
}

main()
