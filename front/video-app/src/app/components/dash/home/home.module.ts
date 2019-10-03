import {NgModule} from '@angular/core';
import {HomeComponent} from './home.component';
import {RouterModule, Routes} from "@angular/router";

const routes: Routes = [
  {
    path: 'dash',
    children: [
      {
        path: '',
        component: HomeComponent,
      }
    ]
  }
];

@NgModule({
  declarations: [
    HomeComponent,
  ],
  imports: [
    RouterModule.forChild(routes),
  ],
  exports: [
    RouterModule
  ]
})
export class HomeModule {
}
