import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BlogCreationComponent } from './blog-creation/blog-creation.component';
import { BlogsListComponent } from './blogs-list/blogs-list.component';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { ReactiveFormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatGridListModule } from '@angular/material/grid-list';

@NgModule({
  declarations: [
    BlogCreationComponent,
    BlogsListComponent
  ],
  imports: [
    CommonModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatChipsModule,
    MatIconModule,
    ReactiveFormsModule,
    MatCardModule,
    MatGridListModule
  ]
})
export class BlogsModule { }