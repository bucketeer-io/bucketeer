import test from 'ava';
import { Bucketeer } from '../../bucketeer';

test('murmur128 should return correct high and low values', (t) => {
  const bucketeer = new Bucketeer();
  const input = 'fid-uid-sampling-seed';

  // Accessing the private method using TypeScript casting
  const { high, low } = (bucketeer as any).murmur128(input);

  t.is(high.toString(), BigInt('2548757552806388169').toString());
  t.is(low.toString(), BigInt('9787172855444729749').toString());
});

test('toFloat64 should return correct normalized value', (t) => {
  const bucketeer = new Bucketeer();
  const high = BigInt('2548757552806388169');
  const low = BigInt('9787172855444729749');

  // Accessing the private method
  const normalized = (bucketeer as any).toFloat64(high, low);

  t.is(normalized, 0.1381684237945762);
});
