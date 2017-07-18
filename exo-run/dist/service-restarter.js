// Generated by LiveScript 1.5.0
var watch, debounce, EventEmitter, DockerCompose, path, ServiceRestarter;
watch = require('chokidar').watch;
debounce = require('debounce');
EventEmitter = require('events').EventEmitter;
DockerCompose = require('../../exosphere-shared').DockerCompose;
path = require('path');
ServiceRestarter = (function(superclass){
  var prototype = extend$((import$(ServiceRestarter, superclass).displayName = 'ServiceRestarter', ServiceRestarter), superclass).prototype, constructor = ServiceRestarter;
  function ServiceRestarter(arg$){
    this.role = arg$.role, this.serviceLocation = arg$.serviceLocation, this.env = arg$.env, this.logger = arg$.logger;
    this.write = bind$(this, 'write', prototype);
    this.watch = bind$(this, 'watch', prototype);
    this.stableDelay = 500;
    this.debounceDelay = 500;
    this.dockerConfigLocation = path.join(process.cwd(), 'tmp');
    this.debouncedRestart = debounce(this._restart, this.debounceDelay);
  }
  ServiceRestarter.prototype.watch = function(){
    /* Ignores any sub-path including dotfiles.
    '[\/\\]' accounts for both windows and unix systems, the '\.' matches a single '.', and the final '.' matches any character. */
    var x$, this$ = this;
    x$ = this.watcher = watch(this.serviceLocation, {
      awaitWriteFinish: {
        stabilityThreshold: this.stableDelay
      },
      ignoreInitial: true,
      ignored: [/.*\/node_modules\/.*/, /(^|[\/\\])\../]
    });
    x$.on('add', function(addedPath){
      this$.logger.log({
        role: 'exo-run',
        text: "Restarting service '" + this$.role + "' because " + addedPath + " was created"
      });
      return this$.debouncedRestart();
    });
    x$.on('change', function(changedPath){
      this$.logger.log({
        role: 'exo-run',
        text: "Restarting service '" + this$.role + "' because " + changedPath + " was changed"
      });
      return this$.debouncedRestart();
    });
    x$.on('unlink', function(removedPath){
      this$.logger.log({
        role: 'exo-run',
        text: "Restarting service '" + this$.role + "' because " + removedPath + " was deleted"
      });
      return this$.debouncedRestart();
    });
    return x$;
  };
  ServiceRestarter.prototype._restart = function(){
    var cwd, this$ = this;
    this.watcher.close();
    cwd = this.dockerConfigLocation;
    return DockerCompose.killContainer({
      serviceName: this.role,
      cwd: cwd,
      write: this.write
    }, function(exitCode){
      switch (false) {
      case !exitCode:
        return this$.emit('error', "Docker failed to kill container " + this$.role);
      }
      this$.write("Docker container stopped");
      return DockerCompose.createNewContainer({
        serviceName: this$.role,
        cwd: cwd,
        env: this$.env,
        write: this$.write
      }, function(exitCode){
        switch (false) {
        case !exitCode:
          return this$.emit('error', "Docker image failed to rebuild " + this$.role);
        }
        this$.write("Docker image rebuilt");
        return DockerCompose.startContainer({
          serviceName: this$.role,
          cwd: cwd,
          env: this$.env,
          write: this$.write
        }, function(exitCode){
          switch (false) {
          case !exitCode:
            return this$.emit('error', "Docker container failed to restart " + this$.role);
          }
          this$.watch();
          return this$.logger.log({
            role: 'exo-run',
            text: "'" + this$.role + "' restarted successfully"
          });
        });
      });
    });
  };
  ServiceRestarter.prototype.write = function(text){
    return this.logger.log({
      role: this.role,
      text: text,
      trim: true
    });
  };
  return ServiceRestarter;
}(EventEmitter));
module.exports = ServiceRestarter;
function bind$(obj, key, target){
  return function(){ return (target || obj)[key].apply(obj, arguments) };
}
function extend$(sub, sup){
  function fun(){} fun.prototype = (sub.superclass = sup).prototype;
  (sub.prototype = new fun).constructor = sub;
  if (typeof sup.extended == 'function') sup.extended(sub);
  return sub;
}
function import$(obj, src){
  var own = {}.hasOwnProperty;
  for (var key in src) if (own.call(src, key)) obj[key] = src[key];
  return obj;
}