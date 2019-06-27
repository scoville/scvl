{{define "body"}}
<div class="jumbotron">
  <h1>SCVL Toolbox</h1>
  <ul class="menu">
    <li>
      <a href="{{.MainHost}}/pages">URL Shortener</a>
    </li>
    <li>
      <a href="{{.FileHost}}">File Uploader</a>
    </li>
    <li>
      <a href="{{.ImageHost}}">Image Uploader</a>
    </li>
  </ul>
</div>
{{end}}
