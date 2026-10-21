/* Theme toggle, inlined in <head> (budget: 1 KB). Applies a stored
   override before first paint; reveals the [data-theme-toggle] button
   (rendered hidden, so no-JS readers never see a dead control); a click
   flips light/dark. Choosing what the system already prefers clears the
   override, so the page follows prefers-color-scheme again. */
(function (d, k) {
var r = d.documentElement, m = matchMedia("(prefers-color-scheme: dark)");
try { var s = localStorage.getItem(k); if (s == "dark" || s == "light") r.dataset.theme = s; } catch (e) {}
d.addEventListener("DOMContentLoaded", function () {
var b = d.querySelector("[data-theme-toggle]");
if (!b) return;
b.hidden = false;
b.addEventListener("click", function () {
var sys = m.matches ? "dark" : "light", next = (r.dataset.theme || sys) == "dark" ? "light" : "dark";
if (next == sys) delete r.dataset.theme; else r.dataset.theme = next;
try { next == sys ? localStorage.removeItem(k) : localStorage.setItem(k, next); } catch (e) {}
});
});
})(document, "theme");
