export interface Outsider{
  file: string
  unmatchedPatterns: Array<string>
}
export interface Violation {
  sha: string
  header: string
  outsiders: Array<Outsider>
}

export interface ConfigData {
  regex: string
  patterns: Record<string, Array<string>>
}
