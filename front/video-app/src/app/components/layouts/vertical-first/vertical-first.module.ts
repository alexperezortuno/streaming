import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {VerticalFirstComponent} from './vertical-first.component';
import {ContentModule} from "../../content/content.module";
import {NavbarModule} from "../../shared/navbar/navbar.module";

@NgModule({
  declarations: [
    VerticalFirstComponent
  ],
  imports: [
    CommonModule,
    ContentModule,
    NavbarModule,
  ],
  exports: [
    VerticalFirstComponent
  ]
})
export class VerticalFirstModule {
}
