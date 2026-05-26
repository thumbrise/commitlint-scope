import { Command, Flags } from '@oclif/core'

export default class Run extends Command {
  static args = {}
  static description = 'Lint commits scopes'
  static examples = [
    `<%= config.bin %> <%= command.id %> --from main --to feature-branch
<%= config.bin %> <%= command.id %> --from HEAD~5 --to HEAD
<%= config.bin %> <%= command.id %> --from $(git merge-base main HEAD) --to HEAD
`,
  ]
  static flags = {
    from: Flags.string({
      description: 'start of commit range (exclusive)',
      helpValue: '<sha>',
      required: true,
    }),
    to: Flags.string({
      description: 'end of commit range (inclusive)',
      helpValue: '<sha>',
      required: true,
    }),
  }

  async run(): Promise<void> {
    const { args, flags } = await this.parse(Run)

    this.log(`hello ${args.person} from ${flags.from}! (./src/commands/hello/index.ts)`)
  }
}
