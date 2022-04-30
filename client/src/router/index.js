import { createRouter, createWebHistory } from "vue-router";
import HelloWorld from "@/components/HelloWorld";
import Signup from "@/components/Signup";
import Signin from "@/components/Signin";

const routes = [
  {
    path: "/:pathMatch(.*)*",
    redirect: "/signin",
  },
  {
    path: "/",
    name: "HelloWorld",
    component: HelloWorld,
  },
  {
    path: "/signup",
    name: "Signup",
    component: Signup,
  },
  {
    path: "/signin",
    name: "Signin",
    component: Signin,
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
