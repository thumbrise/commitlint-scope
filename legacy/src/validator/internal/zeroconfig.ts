
const SCOPE_REGEX = /^[a-z]+(?:\((?<scope>[^)]+)\))?!?:\s/;
function patterns(scope: string): Array<string> {
  return [scope + '/**']
}

export default {
  patterns,
  SCOPE_REGEX,
}
