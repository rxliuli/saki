const os = require('os')
const json = require('./package.json')
const { downloadRelease } = require('@terascope/fetch-github-release')
const path = require('path')
const { extract } = require('tar')
const { copy, remove, pathExists } = require('fs-extra')

async function main() {
  const tempPath = path.resolve(__dirname, '.temp')
  const ext = platform === 'windows' ? '.exe' : ''
  const binPath = path.resolve(__dirname, 'bin' + ext)
  if (await pathExists(binPath)) {
    console.log('已下载')
    return
  }
  const archMap = {
    arm64: 'arm64',
    x64: 'amd64',
  }
  const platformMap = {
    win32: 'windows',
    linux: 'linux',
    darwin: 'macos',
  }
  const platform = platformMap[os.platform()]
  const arch = os.arch()
  const assetName = `saki_${json.version}_${platform}_${archMap[arch]}.tar.gz`
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

  await copy(path.resolve(tempPath, 'saki' + ext), binPath)
  await remove(path.resolve(__dirname, 'bin'))
}

main()
