import Git from "./internal/git.js";
import {OutsiderFinder} from "./internal/outsider-finder.js";
import type {Violation} from "./types.d.ts";
import Config from "./internal/config.js";
import ScopeParser from "./internal/scope-parser.js";
import {LEVEL_DEBUG, Logger} from "./internal/log.js";

export default class Validator {
  private readonly logger: Logger;
  private readonly git: Git;
  private readonly outsiderFinder: OutsiderFinder;
  private readonly scopeParser: ScopeParser;

  public constructor() {
    this.logger = new Logger(LEVEL_DEBUG)
    this.git = new Git()
    const config = new Config()
    this.outsiderFinder = new OutsiderFinder(config)
    this.scopeParser = new ScopeParser(config)
  }

  public async validate(from: string, to: string): Promise<Array<Violation>> {
    const shaRange = await this.git.shaRange(from, to)
    const violations: Array<Violation> = []
    for (const sha of shaRange) {
      const message = await this.git.message(sha)
      if (!message) {
        this.logger.debug('No message. Skip.', {sha})

        continue
      }
      const scope = this.scopeParser.parse(message)
      if (!scope) {
        this.logger.debug('No scope. Skip.', {sha, message})

        continue
      }

      const files = await this.git.filesChanged(sha)
      if (files.length === 0) {
        this.logger.debug('No files changed. Skip.', {sha, message})

        continue
      }

      const outsiders = this.outsiderFinder.find(scope, files)
      if (outsiders.length > 0) {
        violations.push({
          sha: sha.slice(0, 7),
          header: message,
          outsiders: outsiders,
        })
      }
    }

    return violations
  }

}

