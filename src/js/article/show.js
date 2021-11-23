'use struct';

document.addEventListener('DOMContentLoaded', function() {
  const elm = document.getElementById('article-body');
  elm.innerHTML = md.render(elm.dataset.markdown);
})