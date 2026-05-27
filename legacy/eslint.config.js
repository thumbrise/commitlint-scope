import { includeIgnoreFile } from '@eslint/compat'
import oclif from 'eslint-config-oclif'
import prettier from 'eslint-config-prettier'
import { defineConfig, globalIgnores } from 'eslint/config'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

function file(name) {
  return path.resolve(path.dirname(fileURLToPath(import.meta.url)), name)
}

export default defineConfig([
  includeIgnoreFile(file('.gitignore')),
  globalIgnores(['.releaserc.js']),
  ...oclif,
  prettier,
])
