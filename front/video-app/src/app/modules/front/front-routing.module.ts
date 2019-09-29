import {NgModule} from '@angular/core';
import {RouterModule, Routes} from "@angular/router";

const routes: Routes = [
    {
        path: '',
        loadChildren: () => import('../../components/front/login/login.module')
        .then(m => m.LoginModule)
    },
    {
        path: '',
        loadChildren: () => import('../../components/front/dash/dash.module')
        .then(m => m.DashModule )
    }
];

@NgModule({
    imports: [
        RouterModule.forChild(routes)
    ],
    exports: [
        RouterModule
    ]
})
export class FrontRoutingModule {
}
