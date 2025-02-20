import murmurHash3 from 'murmurhash3js';

export default class Bucketeer {
  // Calculate the input hash of the target property and map it to a float64 between [0,1]
  bucket(input: string): number {
    const { high, low } = this.murmur128(input);
    return this.toFloat64(high, low);
  }

  private murmur128(input: string): { high: bigint; low: bigint } {
    // Compute the MurmurHash3 (128-bit) hash using the native library.
    // By default, murmurHash3.x64.hash128 returns the hash in big-endian format,
    // which is crucial for consistency with the Go implementation that also uses big-endian.
    const hash = murmurHash3.x64.hash128(input);

    // Convert the hash from hexadecimal to BigInt.
    // The first 16 characters represent the high part (most significant bits),
    // and the last 16 characters represent the low part (least significant bits).
    const high = BigInt('0x' + hash.slice(0, 16)); // High part
    const low = BigInt('0x' + hash.slice(16, 32)); // Low part

    return { high, low };
  }

  private toFloat64(high: bigint, low: bigint): number {
    // Combine the high and low parts into a single floating-point number.
    // This maintains the full 128-bit range, ensuring a normalized value between [0,1].
    const full = Number(high) * Math.pow(2, 64) + Number(low);
    const maxValue = Math.pow(2, 128) - 1; // Maximum value for 128-bit

    return full / maxValue; // Normalize to the range [0,1]
  }
}

export { Bucketeer };
