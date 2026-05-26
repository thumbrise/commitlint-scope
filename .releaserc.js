/*
 * Copyright 2026 thumbrise
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

const fs = require('node:fs')
const path = require('node:path')

const commitPartial = fs.readFileSync(path.join(__dirname, 'release-template.hbs'), 'utf8')

module.exports = {
  branches: ['main'],
  plugins: [
    [
      '@semantic-release/commit-analyzer',
      {
        preset: 'conventionalcommits',
      },
    ],
    [
      '@semantic-release/release-notes-generator',
      {
        parserOpts: {
          noteKeywords: ['BREAKING CHANGE', 'BREAKING CHANGES', 'BREAKING', '!'],
        },
        preset: 'conventionalcommits',
        presetConfig: {
          types: [
            { section: 'Features', type: 'feat' },
            { section: 'Bug Fixes', type: 'fix' },
            { section: 'CI/CD', type: 'ci' },
            { section: 'Tests', type: 'test' },
            { section: 'Reverts', type: 'revert' },
            { section: 'Build System', type: 'build' },
            { section: 'Code Refactoring', type: 'refactor' },
            { section: 'Code Refactoring', type: 'style' },
            { section: 'Performance Improvements', type: 'perf' },
            { section: 'Documentation', type: 'docs' },
            { section: 'Internal Changes', type: 'chore' },
          ],
        },
        writerOpts: {
          bodyWrap: 100,
          commitPartial,
          commitsSort: ['scope', 'subject'],
          includeDetails: true,
          showBody: true,
        },
      },
    ],
    [
      '@semantic-release/exec',
      {
        prepareCmd: 'npm run version',
      },
    ],
    [
      '@semantic-release/npm',
      {
        npmPublish: false,
      },
    ],
    [
      '@semantic-release/git',
      {
        assets: ['package.json', 'README.md'],
        message: 'chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}',
      },
    ],
    '@semantic-release/github',
  ],
}
