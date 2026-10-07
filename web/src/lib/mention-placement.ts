// Placement carried over from the last approved prototype. Coordinates use visualViewport.
export function mentionPlacement(
  anchor: { left: number; width: number },
  caret: { top: number; height: number },
  viewport: { left: number; top: number; width: number; height: number },
  headerBottom = 0,
  contentHeight = 240,
) {
  const leftEdge = viewport.left + 8,
    rightEdge = viewport.left + viewport.width - 8;
  const topEdge = viewport.top + 8,
    bottomEdge = viewport.top + viewport.height - 8;
  const topLimit =
    viewport.height < 280
      ? topEdge
      : Math.min(bottomEdge - 64, Math.max(topEdge, headerBottom + 6));
  const width = Math.min(360, Math.max(260, anchor.width), rightEdge - leftEdge);
  const left = Math.max(leftEdge, Math.min(anchor.left, rightEdge - width));
  const desired = Math.max(40, Math.min(240, contentHeight));
  const caretTop = Math.max(topLimit, Math.min(caret.top, bottomEdge - caret.height));
  const belowStart = caretTop + caret.height + 6;
  const below = Math.max(0, bottomEdge - belowStart),
    above = Math.max(0, caretTop - 6 - topLimit);
  const down = below >= Math.min(150, desired) || below >= above;
  const height = Math.max(40, Math.min(desired, down ? below : above, bottomEdge - topLimit));
  const top = Math.max(
    topLimit,
    Math.min(down ? belowStart : caretTop - 6 - height, bottomEdge - height),
  );
  return { left, top, width, height };
}
