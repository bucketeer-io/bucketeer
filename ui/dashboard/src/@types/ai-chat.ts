export type AIChatRole = 'user' | 'assistant';

export interface AIChatMessage {
  id: string;
  role: AIChatRole;
  content: string;
}

export type PageType =
  | 'feature_flags'
  | 'targeting'
  | 'experiments'
  | 'segments'
  | 'autoops';

export interface PageContext {
  pageType: PageType | '';
  featureId?: string;
  metadata?: Record<string, string>;
}

export type SuggestionType =
  | 'SUGGESTION_TYPE_UNSPECIFIED'
  | 'SUGGESTION_TYPE_FEATURE_DISCOVERY'
  | 'SUGGESTION_TYPE_BEST_PRACTICE'
  | 'SUGGESTION_TYPE_OPTIMIZATION'
  | 'SUGGESTION_TYPE_WARNING';

export interface Suggestion {
  id: string;
  type: SuggestionType;
  title: string;
  description: string;
  docUrl: string;
}
