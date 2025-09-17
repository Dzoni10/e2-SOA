import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TourExecutionComponent } from './tour-execution.component';

describe('TourExecutionComponent', () => {
  let component: TourExecutionComponent;
  let fixture: ComponentFixture<TourExecutionComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [TourExecutionComponent]
    });
    fixture = TestBed.createComponent(TourExecutionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
