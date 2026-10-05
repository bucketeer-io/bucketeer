export interface UnsavedForm {
  isShow: boolean;
  title: string;
  content: string;
  callBackCancel?: () => void;
}

// Keep references so navigation always reads the latest state and callbacks.
export function createUnsavedForms() {
  const forms = new Set<{ current: UnsavedForm }>();
  return {
    register(form: { current: UnsavedForm }) {
      forms.add(form);
      return () => {
        forms.delete(form);
      };
    },
    getDirtyForms: () =>
      Array.from(forms, form => form.current).filter(form => form.isShow)
  };
}
