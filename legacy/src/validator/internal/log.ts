export const LEVEL_TRACE = 10
export const LEVEL_DEBUG = 20
export const LEVEL_INFO  = 30
export const LEVEL_WARN  = 40
export const LEVEL_ERROR = 50
export const LEVEL_FATAL = 60
interface LogRecord {
  level: number
  message: string
  context?: Record<string, unknown>
}
export class Logger {
  public constructor(private readonly minLevel: number) {
  }
  public log(record: LogRecord): void {
    if (record.level < this.minLevel) return

    const timestamp = new Date().toISOString()
    const levelName = this.getLevelName(record.level)
    const ctx = record.context ? ` ${JSON.stringify(record.context)}` : ''
    const formatted = `[${timestamp}] [${levelName}] ${record.message}${ctx}`

    // Используем соответствующий метод консоли, иначе console.log
    switch (record.level) {
      case LEVEL_TRACE:
      case LEVEL_DEBUG:
        console.debug(formatted)
        break
      case LEVEL_INFO:
        console.info(formatted)
        break
      case LEVEL_WARN:
        console.warn(formatted)
        break
      case LEVEL_ERROR:
      case LEVEL_FATAL:
        console.error(formatted)
        break
      default:
        console.log(formatted)
    }
  }

  // Хелперы
  public trace(message: string, context?: Record<string, unknown>): void {
    this.log({ level: LEVEL_TRACE, message, context })
  }

  public debug(message: string, context?: Record<string, unknown>): void {
    this.log({ level: LEVEL_DEBUG, message, context })
  }

  public info(message: string, context?: Record<string, unknown>): void {
    this.log({ level: LEVEL_INFO, message, context })
  }

  public warn(message: string, context?: Record<string, unknown>): void {
    this.log({ level: LEVEL_WARN, message, context })
  }

  public error(message: string, error?: Error, context?: Record<string, unknown>): void {
    // Error передаём как часть контекста
    const errorCtx = error ? { error: error.message, stack: error.stack, ...context } : context
    this.log({ level: LEVEL_ERROR, message, context: errorCtx })
  }

  public fatal(message: string, error?: Error, context?: Record<string, unknown>): void {
    const errorCtx = error ? { error: error.message, stack: error.stack, ...context } : context
    this.log({ level: LEVEL_FATAL, message, context: errorCtx })
  }

  private getLevelName(level: number): string {
    switch (level) {
      case LEVEL_TRACE: return 'TRACE'
      case LEVEL_DEBUG: return 'DEBUG'
      case LEVEL_INFO:  return 'INFO'
      case LEVEL_WARN:  return 'WARN'
      case LEVEL_ERROR: return 'ERROR'
      case LEVEL_FATAL: return 'FATAL'
      default:          return `LEVEL(${level})`
    }
  }
}
