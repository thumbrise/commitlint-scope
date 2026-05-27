import {cosmiconfigSync} from 'cosmiconfig';
import zeroconfig from "./zeroconfig.js";

import type {ConfigData} from "../types.d.ts";

// For example .commitlint-scoperc.js
const MODULE_NAME = 'commitlint-scope';

export default class Config {
  private cfg: ConfigData | null = null

  public constructor() {
    const searchResult = cosmiconfigSync(MODULE_NAME).search()
    this.cfg = searchResult?.config || null
  }

  public scopePatterns(scope: string): Array<string> {
    return this.cfg?.patterns[scope] || zeroconfig.patterns(scope)
  }

  public scopeRegex(): RegExp {
    let regex = this.cfg?.regex;
    if (!regex) {
      return zeroconfig.SCOPE_REGEX
    }

    return new RegExp(regex)
  }

}


