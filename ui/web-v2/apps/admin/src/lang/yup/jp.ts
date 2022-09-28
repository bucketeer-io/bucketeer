export const localJp = {
  mixed: {
    default: '入力エラーです',
    required: '必須入力項目です',
    oneOf: ({ values }) => `次の値のいずれかでなければなりません: ${values}`,
    notOneOf: ({ values }) => `次の値のいずれかであってはなりません: ${values}`,
    isNumber: '形式が違います',
  },
  string: {
    required: '必須入力項目です',
    length: ({ length }) => `${length}文字入力して下さい`,
    min: ({ min }) => `${min}文字以上入力してください`,
    max: ({ max }) => `${max}文字以内で入力して下さい`,
    matches: '形式が違います',
    email: '形式が違います',
    url: '形式が違います',
    trim: '前後にスペースを入れてはいけません',
    lowercase: '小文字でなければなりません',
    uppercase: '大文字でなければなりません',
  },
  number: {
    min: ({ min }) => `${min}以上の値を入力して下さい`,
    max: ({ max }) => `${max}以下の値を入力して下さい`,
    lessThan: ({ less }) => `${less}より小さい値を入力して下さい`,
    moreThan: ({ more }) => `${more}より大きい値を入力して下さい`,
    notEqual: ({ notEqual }) => `${notEqual}と違う値を入力して下さい`,
    positive: '正の数を入力して下さい',
    negative: '負の数を入力して下さい',
    integer: '整数を入力して下さい',
  },
  date: {
    default: '形式が違います',
    min: ({ min }) => `${min}以上の日付を入力して下さい`,
    max: ({ max }) => `${max}以下の日付を入力して下さい`,
  },
  object: {
    noUnknown: '有効なキーを持ったデータを入力して下さい',
  },
  array: {
    min: ({ min }) => `${min}個以上の値を入力して下さい`,
    max: ({ max }) => `${max}個以下の値を入力して下さい`,
  },
};
