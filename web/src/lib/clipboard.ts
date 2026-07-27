export async function copyText(value: string): Promise<boolean> {
  if (window.isSecureContext && navigator.clipboard) {
    const wroteViaClipboardApi = await navigator.clipboard.writeText(value).then(
      () => true,
      () => false,
    );
    if (wroteViaClipboardApi) return true;
  }

  try {
    const textarea = document.createElement("textarea");
    textarea.value = value;
    textarea.style.position = "fixed";
    textarea.style.opacity = "0";
    document.body.appendChild(textarea);
    textarea.focus();
    textarea.select();
    const ok = document.execCommand("copy");
    document.body.removeChild(textarea);
    return ok;
  } catch {
    return false;
  }
}
