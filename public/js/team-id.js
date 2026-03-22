/**
 * Normalize UUID strings for comparison.
 * PostgREST and localStorage can differ in casing; strict === misses and breaks active-team selection.
 */
export function normId(id) {
    if (id == null || id === '') return '';
    return String(id).trim().toLowerCase();
}
