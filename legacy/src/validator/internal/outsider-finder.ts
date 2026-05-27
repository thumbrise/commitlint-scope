import micromatch from 'micromatch'
import Config from "./config.js";
import type {Outsider} from "../types.d.ts";


export class OutsiderFinder {
  public constructor(private readonly config: Config) {
  }

  public find(scope: string, files: Array<string>): Array<Outsider> {
    let patterns = this.config.scopePatterns(scope);
    let result: Array<Outsider> = []
    files.forEach((file) => {
      let violated = !micromatch.isMatch(file, patterns, {dot: true})
      if (violated) {
        result.push({
          file: file,
          unmatchedPatterns: patterns,
        })
      }
    })
    return result
  }

}
