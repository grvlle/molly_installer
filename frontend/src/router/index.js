import Vue from "vue";
import Router from "vue-router";
import Home from "@/views/Home";
import Install from "@/views/Install";
import Uninstall from "@/views/Uninstall";

Vue.use(Router);

const router = new Router({
  mode: "abstract",
  routes: [
    {
      path: "/",
      name: "Home",
      component: Home,
    },
    {
      path: "/install",
      name: "Install",
      component: Install,
    },
    {
      path: "/uninstall",
      name: "Uninstall",
      component: Uninstall,
    },
  ],
});

export default router;

/**
 * Asynchronously load view (Webpack Lazy loading compatible)
 * The specified component must be inside the Views folder
 * @param  {string} name  the filename (basename) of the view to load.
function view(name) {
   var res= require('../components/Dashboard/Views/' + name + '.vue');
   return res;
};**/
