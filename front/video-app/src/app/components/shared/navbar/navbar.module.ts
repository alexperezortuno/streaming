import {NgModule} from '@angular/core';
import {NavbarComponent} from "./navbar.component";
import {MatToolbarModule} from "@angular/material/toolbar";
import {MatIconModule} from "@angular/material/icon";
import {MatButtonModule} from "@angular/material/button";
import {RouterModule} from "@angular/router";

@NgModule({
  declarations: [
    NavbarComponent
  ],
  imports: [
    MatToolbarModule,
    MatIconModule,
    MatButtonModule,
    RouterModule
  ],
  exports: [
    NavbarComponent
  ]
})
export class NavbarModule {
}
