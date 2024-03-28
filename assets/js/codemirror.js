console.log('entra')
require([
  "./js/lib/codemirror/lib/codemirror",
  "./js/lib/codemirror/mode/sql/sql",
  "./js/lib/codemirror/addon/hint/show-hint",
  "./js/lib/codemirror/addon/hint/sql-hint",
], function(CodeMirror) {
  const attachEditor = () => {
    document.querySelectorAll('.code-editor').forEach((editorElement) => {
      const content = editorElement.textContent;
      editorElement.innerHTML = '';
      var editor = CodeMirror(editorElement, {
        extraKeys: { "Ctrl-Space": "autocomplete" },
        value: content.trim(),
        lineNumbers: true,
        lineWrapping: true,
        mode: "text/x-sql",
        theme: "darcula",
      });

      console.log('editor created', editor)
    });
  }
  attachEditor();
  htmx.onLoad(function() {
    attachEditor();
  });

});
