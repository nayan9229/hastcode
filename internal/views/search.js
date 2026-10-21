/* Search page loader (page_search.templ). Nothing loads until the Reader
   focuses the input or submits the form (or arrives with ?q=): then the
   Pagefind modular UI is fetched from the form's data-pagefind-bundle,
   which in turn loads the index on demand. Results are styled by site.css,
   not Pagefind's stylesheet. */
(function () {
  var f = document.querySelector("form[data-pagefind-bundle]");
  if (!f || !window.Promise) return;
  var b = f.getAttribute("data-pagefind-bundle"), i = f.querySelector("input"), st = document.getElementById("search-status"), ui;
  function load() {
    if (!ui) ui = new Promise(function (ok, no) {
      var s = document.createElement("script");
      s.src = b + "pagefind-modular-ui.js";
      s.onload = ok;
      s.onerror = no;
      document.head.appendChild(s);
    }).then(function () {
      var M = window.PagefindModularUI, x = new M.Instance({ bundlePath: b });
      x.add(new M.Input({ inputElement: "#" + i.id }));
      x.add(new M.Summary({ containerElement: "#search-status" }));
      x.add(new M.ResultList({ containerElement: "#search-results", showImages: false }));
      x.triggerLoad();
      return x;
    }).catch(function () {
      st.textContent = "Search is unavailable right now. Browse by topic below.";
    });
    return ui;
  }
  /* Search for whatever the input holds once loaded, so text typed while
     Pagefind was still loading is not lost. */
  function go() {
    load().then(function (x) { if (x && i.value) x.triggerSearch(i.value); });
  }
  i.addEventListener("focus", go, { once: true });
  f.addEventListener("submit", function (e) { e.preventDefault(); go(); });
  var q = new URLSearchParams(location.search).get("q");
  if (q) { i.value = q; go(); }
})();
