import {NgModule} from '@angular/core';
import {RouterModule, Routes} from "@angular/router";
import {HomeModule} from "../components/dash/home/home.module";

@NgModule({
  imports: [
    HomeModule
  ],
})
export class DashModule {
}
