import Config from "./config.js";

export default class ScopeParser {
  public constructor(private readonly config: Config) {
  }

  parse(message: string): string | undefined {
    const match = message.match(this.config.scopeRegex());
    return match?.groups?.scope
  }
}
