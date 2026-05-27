import {SimpleGit, simpleGit} from 'simple-git'

export default class Git {
  private readonly git: SimpleGit

  public constructor(baseDir?: string) {
    this.git = simpleGit(baseDir)
  }

  public async filesChanged(sha: string): Promise<Array<string>> {
    const raw = await this.git.raw(['diff-tree', '--no-commit-id', '-r', '--name-only', sha])
    return raw.trim().split('\n').filter(Boolean)
  }

  public async shaRange(from: string, to: string): Promise<Array<string>> {
    let r = await this.git.raw(['rev-list', `${from}..${to}`])
    return r.trim().split('\n').filter(Boolean)
  }


  public async message(sha: string): Promise<string> {
    return this.git.raw(['log', '--format=%s', '-1', sha])
  }
}

