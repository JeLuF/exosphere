require! {
  'fs-extra' : fs
  'path'
  'tmplconv'
}


World = !->

  @create-empty-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    data =
      'app-name': app-name
      'app-description': 'Empty test application'
      'app-version': '1.0.0'
    src-path = path.join process.cwd!, 'node_modules' 'exosphere-shared' 'templates', 'create-app'
    tmplconv.render(src-path, app-dir, {data})


module.exports = ->
  @World = World
