export type ErrorValue<T> = [T, null] | [null, Error];
export type ErrorValuePromize<T> = Promise<ErrorValue<T>>;
export async function errorValue<T>(promise: Promise<T>): ErrorValuePromize<T> {
  try {
    const value = await promise;
    return [value, null];
  } catch (error) {
    if (error instanceof Error) return [null, error];
    return [null, new Error(String(error))];
  }
}
