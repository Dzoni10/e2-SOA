import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BlogCreationComponent } from './blog-creation/blog-creation.component';
import { BlogsListComponent } from './blogs-list/blogs-list.component';
import { BlogCommentsComponent } from './blog-comments/blog-comments.component';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatGridListModule } from '@angular/material/grid-list';
import { MarkdownModule } from 'ngx-markdown';
import { RouterModule } from '@angular/router';
import { BlogEditComponent } from './blog-edit/blog-edit.component'; // DODANO

@NgModule({
  declarations: [
    BlogCreationComponent,
    BlogsListComponent,
    BlogCommentsComponent,
    BlogEditComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    RouterModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatChipsModule,
    MatIconModule,
    ReactiveFormsModule,
    MatCardModule,
    MatGridListModule,
    MarkdownModule.forRoot()

  ]
})
export class BlogsModule { }