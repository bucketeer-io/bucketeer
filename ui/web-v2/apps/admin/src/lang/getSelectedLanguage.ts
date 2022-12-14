export enum LanguageTypes {
  JAPAN = 'ja',
  ENGLISH = 'en',
}

export const getSelectedLanguage = () => {
  const language = window.localStorage.getItem('language');

  if (language) {
    return language;
  }

  const supportedLanguages = [LanguageTypes.JAPAN, LanguageTypes.ENGLISH];

  const foundLanguage = supportedLanguages.find(
    (lang) => lang === window.navigator.language?.slice(0, 2)
  );

  let selectedLanguage;

  if (foundLanguage) {
    selectedLanguage = foundLanguage;
  } else {
    // Default to English if no supported language is found
    selectedLanguage = LanguageTypes.ENGLISH;
  }

  window.localStorage.setItem('language', selectedLanguage);

  return selectedLanguage;
};
