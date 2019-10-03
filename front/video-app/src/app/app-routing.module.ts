import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {LoginModule} from "./components/auth/login/login.module";
import {AuthModule} from "./modules/auth.module";
import {DashModule} from "./modules/dash.module";

const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'auth/login'
  }
];

@NgModule({
  imports: [
    AuthModule,
    DashModule,
    RouterModule.forRoot(
      routes,
      {enableTracing: false} // <-- debugging purposes only
    )
  ],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
