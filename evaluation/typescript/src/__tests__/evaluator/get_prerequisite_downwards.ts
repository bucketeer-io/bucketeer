import test from 'ava';
import {
  createPrerequisite,
  creatFeature,
} from '../../modelFactory';
import { Evaluator } from '../../evaluation';

var allFeaturesForPrerequisiteTest = {
  featureA: creatFeature({
    id: 'featureA',
    name: 'featureA',
    prerequisitesList: [createPrerequisite('featureE', ''), createPrerequisite('featureF', '')],
  }),
  featureB: creatFeature({
    id: 'featureB',
    name: 'featureB',
    prerequisitesList: [],
  }),
  featureC: creatFeature({
    id: 'featureC',
    name: 'featureC',
    prerequisitesList: [createPrerequisite('featureL', '')],
  }),
  featureD: creatFeature({
    id: 'featureD',
    name: 'featureD',
    prerequisitesList: [],
  }),
  featureE: creatFeature({
    id: 'featureE',
    name: 'featureE',
    prerequisitesList: [createPrerequisite('featureG', '')],
  }),
  featureF: creatFeature({
    id: 'featureF',
    name: 'featureF',
    prerequisitesList: [],
  }),
  featureG: creatFeature({
    id: 'featureG',
    name: 'featureG',
    prerequisitesList: [createPrerequisite('featureH', '')],
  }),
  featureH: creatFeature({
    id: 'featureH',
    name: 'featureH',
    prerequisitesList: [createPrerequisite('featureI', ''), createPrerequisite('featureJ', '')],
  }),
  featureI: creatFeature({
    id: 'featureI',
    name: 'featureI',
    prerequisitesList: [createPrerequisite('featureK', '')],
  }),
  featureJ: creatFeature({
    id: 'featureJ',
    name: 'featureJ',
    prerequisitesList: [],
  }),
  featureK: creatFeature({
    id: 'featureK',
    name: 'featureK',
    prerequisitesList: [],
  }),
  featureL: creatFeature({
    id: 'featureL',
    name: 'featureL',
    prerequisitesList: [createPrerequisite('featureM', ''), createPrerequisite('featureN', '')],
  }),
  featureM: creatFeature({
    id: 'featureM',
    name: 'featureM',
    prerequisitesList: [],
  }),
  featureN: creatFeature({
    id: 'featureN',
    name: 'featureN',
    prerequisitesList: [],
  }),
};
/*
graph TD;
  featureA["Feature A"] --> featureE["Feature E"];
  featureA --> featureF["Feature F"];
  
  featureC["Feature C"] --> featureL["Feature L"];
  
  featureE --> featureG["Feature G"];
  
  featureG["Feature G"] --> featureH["Feature H"];
  
  featureH --> featureI["Feature I"];
  featureH --> featureJ["Feature J"];
  
  featureI["Feature I"] --> featureK["Feature K"];
  
  featureL --> featureM["Feature M"];
  featureL --> featureN["Feature N"];
  
  featureB["Feature B"] --- NoPrerequisite;
  featureD["Feature D"] --- NoPrerequisite;
  featureF --- NoPrerequisite;
  featureJ --- NoPrerequisite;
  featureK --- NoPrerequisite;
  featureM --- NoPrerequisite;
  featureN --- NoPrerequisite;
*/
// The order of the output is a bit diffence with golang test, but the result is the same
const TestCases = [
  {
    desc: 'success: No prerequisites',
    target: [
      allFeaturesForPrerequisiteTest.featureB,
      allFeaturesForPrerequisiteTest.featureD,
    ],
    expected: [
      allFeaturesForPrerequisiteTest.featureB,
      allFeaturesForPrerequisiteTest.featureD,
    ],
    expectedErr: null,
  },
  {
    desc: 'success: Get prerequisites pattern1',
    target: [
      allFeaturesForPrerequisiteTest.featureA,
    ],
    expected: [
      allFeaturesForPrerequisiteTest.featureA,
      allFeaturesForPrerequisiteTest.featureE,
      allFeaturesForPrerequisiteTest.featureG,
      allFeaturesForPrerequisiteTest.featureH,
      allFeaturesForPrerequisiteTest.featureI,
      allFeaturesForPrerequisiteTest.featureK,
      allFeaturesForPrerequisiteTest.featureJ,
      allFeaturesForPrerequisiteTest.featureF,
    ],
    expectedErr: null,
  },
  {
    desc: 'success: Get prerequisites pattern2',
    target: [
      allFeaturesForPrerequisiteTest.featureC,
      allFeaturesForPrerequisiteTest.featureD,
    ],
    expected: [
      allFeaturesForPrerequisiteTest.featureC,
      allFeaturesForPrerequisiteTest.featureL,
      allFeaturesForPrerequisiteTest.featureM,
      allFeaturesForPrerequisiteTest.featureN,
      allFeaturesForPrerequisiteTest.featureD,
    ],
    expectedErr: null,
  },
  {
    desc: 'success: Get prerequisites pattern3',
    target: [
      allFeaturesForPrerequisiteTest.featureD,
      allFeaturesForPrerequisiteTest.featureH,
    ],
    expected: [
      allFeaturesForPrerequisiteTest.featureD,
      allFeaturesForPrerequisiteTest.featureH,
      allFeaturesForPrerequisiteTest.featureI,
      allFeaturesForPrerequisiteTest.featureK,
      allFeaturesForPrerequisiteTest.featureJ,

    ],
    expectedErr: null,
  },
];

const allFeatures = Object.values(allFeaturesForPrerequisiteTest);

TestCases.forEach(({ desc, target, expected, expectedErr }) => {
  test(desc, async (t) => {
    try {
      const evalator = new Evaluator();
      const actual = evalator.getPrerequisiteDownwards(target, allFeatures);
      t.deepEqual(actual, expected);
    } catch (error) {
      t.deepEqual(error, expectedErr);
    }
  });
});