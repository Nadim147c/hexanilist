export type Result<T> = [T, null] | [null, Error]
export type AsyncResult<T> = Promise<Result<T>>

export function handle<T>(fun: () => T): Result<T> {
  try {
    const value = fun()
    return [value, null]
  } catch (error) {
    if (error instanceof Error) return [null, error]
    return [null, new Error(String(error))]
  }
}

export async function handlePromise<T>(promise: Promise<T>): AsyncResult<T> {
  try {
    const value = await promise
    return [value, null]
  } catch (error) {
    if (error instanceof Error) return [null, error]
    return [null, new Error(String(error))]
  }
}
